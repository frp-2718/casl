package bib

var mwxmls = [][]byte{
	[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="multiwhere">
  <query>
    <ppn>ppn000001</ppn>
    <result>
      <library>
        <rcr>rcr000001</rcr>
        <shortname>TEST1</shortname>
      </library>
      <library>
        <rcr>rcr000002</rcr>
        <shortname>TEST2</shortname>
      </library>
      <library>
        <rcr>rcr000003</rcr>
        <shortname>TEST3</shortname>
      </library>
    </result>
  </query>
  <query>
    <ppn>ppn000002</ppn>
    <result>
      <library>
        <rcr>rcr000001</rcr>
        <shortname>TEST1</shortname>
      </library>
    </result>
  </query>
  </sudoc>`),
	[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="multiwhere">
  <query>
    <ppn>ppn000003</ppn>
    <result>
      <library>
        <rcr>rcr000001</rcr>
        <shortname>TEST1</shortname>
      </library>
    </result>
  </query>
  </sudoc>`)}

var mwsemi = [][]byte{[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="multiwhere">
  <query>
    <ppn>ppn000001</ppn>
    <result>
      <library>
        <rcr>rcr000001</rcr>
        <shortname>TEST1</shortname>
      </library>
      <library>
        <rcr>rcr000002</rcr>
        <shortname>TEST2</shortname>
      </library>
      <library>
        <rcr>rcr000003</rcr>
        <shortname>TEST3</shortname>
      </library>
    </result>
  </query>
  </sudoc>`)}

var almaNoLoc = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <bibs total_record_count="0"/>`)

var almaErrMMS = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <web_service_result xmlns="">
    <errorsExist>true</errorsExist>
    <errorList>
      <error>
        <errorCode>402204</errorCode>
        <errorMessage>Input parameters mmsId mmsa is not numeric.</errorMessage>
        <trackingId></trackingId>
      </error>
    </errorList>
  </web_service_result>`)

var almaErrSys = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <web_service_result xmlns="">   
    <errorsExist>true</errorsExist>
    <errorList>
      <error>
        <errorCode>401652</errorCode>
        <errorMessage>Input parameters mmsId mmsa is not numeric.</errorMessage>
        <trackingId></trackingId>
      </error>
    </errorList>
  </web_service_result>`)

var almaOk = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <bibs total_record_count="1">
    <bib>
      <mms_id>id1</mms_id>
      <record_format>unimarc</record_format>
      <linked_record_id/>
      <title>Microéconomie</title>
      <author>Etner , François</author>
      <isbn>2-13-043191-7</isbn>
      <network_numbers>
        <network_number>L/T</network_number>
        <network_number>L/T 11</network_number>
        <network_number>lutest</network_number>
        <network_number>(PPN)123456789</network_number>
      </network_numbers>
      <place_of_publication>Paris</place_of_publication>
      <date_of_publication>DL 1991</date_of_publication>
      <publisher_const>Presses universitaires de France</publisher_const>
      <created_date>2020-11-05Z</created_date>
      <brief_level desc="10">10</brief_level>
      <record>
        <leader>     cam2 2200685   450 </leader>
        <controlfield tag="001">ctfield001</controlfield>
        <controlfield tag="005">ctfield005</controlfield>
        <datafield ind1=" " ind2=" " tag="010">
          <subfield code="a">2-13-043191-7</subfield>
          <subfield code="b">br.</subfield>
          <subfield code="d">125 FRF</subfield>
        </datafield>
        <datafield ind1=" " ind2=" " tag="AVA">
          <subfield code="0">code</subfield>
          <subfield code="8">code</subfield>
          <subfield code="b">BIB1</subfield>
          <subfield code="c">localisation1</subfield>
          <subfield code="d">888</subfield>
          <subfield code="e">available</subfield>
          <subfield code="f">1</subfield>
          <subfield code="g">0</subfield>
          <subfield code="j">LOC1</subfield>
          <subfield code="q">localisation</subfield>
        </datafield>
        <datafield ind1=" " ind2=" " tag="AVA">
          <subfield code="0">999</subfield>
          <subfield code="8">222</subfield>
          <subfield code="b">BIB2</subfield>
          <subfield code="d">g 1</subfield>
          <subfield code="e">available</subfield>
          <subfield code="f">1</subfield>
          <subfield code="g">0</subfield>
          <subfield code="j">LOC2</subfield>
          <subfield code="q">localisation</subfield>
        </datafield>
        <datafield ind1=" " ind2=" " tag="SYS">
          <subfield code="a">000</subfield>
        </datafield>
      </record>
    </bib>
  </bibs>`)

var sudocOk = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="multiwhere">
<query>
<ppn>PPN1</ppn>
<result>
<library>
<rcr>RCR1</rcr>
<shortname>BIB1</shortname>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
<library>
<rcr>RCR_IGN1</rcr>
<shortname>BIB2</shortname>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
</result>
</query>
<query>
<ppn>PPN2</ppn>
<result>
<library>
<rcr>RCR2</rcr>
<shortname>BIB2</shortname>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
<library>
<rcr>RCR_IGN2</rcr>
<shortname>BIB3</shortname>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
<library>
<rcr>RCR3</rcr>
<shortname>BIB4</shortname>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
</result>
</query>
</sudoc>`)

var sudocUnknown = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="multiwhere">
<error>Found a null xml in result : values={ppn=00xx00}, query=select autorites.MULTIWHERE(#ppn#) from dual </error>
</sudoc>`)

var rcrUnknown = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="iln2rcr"><error>found a null xml in result : values={iln=a}, query=select autorites.iln2rcr(#iln#) from dual </error></sudoc>`)

var rcrOk = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="iln2rcr">
	<query>
		<iln>4</iln>
		<result>
			<library>
				<rcr>000000001</rcr>
				<name>UNIV1</name>
				<shortname>BIB1</shortname>
				<oclcsymbol></oclcsymbol>
<address>ADRESSE1</address>
<tel>TEL1</tel>
<web>WEB1</web>
<type>bibliothèque</type>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
<library>
<rcr>000000002</rcr>
<name>UNIV2</name>
<shortname>BIB2</shortname>
<oclcsymbol></oclcsymbol>
<address>ADRESSE2</address>
<tel>TEL2</tel>
<web>WEB2</web>
<type>Bibliothèque</type>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
</result>
</query>
<query>
<iln>4</iln>
<result>
<library>
<rcr>000000003</rcr>
<name>UNIV3</name>
<shortname>BIB3</shortname>
<oclcsymbol></oclcsymbol>
<address>ADRESSE3</address>
<tel>TEL3</tel>
<web>WEB3</web>
<type>Bibliothèque</type>
<latitude>99.9</latitude>
<longitude>99.9</longitude>
</library>
</result>
</query>
</sudoc>`)

var fakeAlmaRes = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
  <bibs total_record_count="1">
    <bib>
      <record>
        <leader>     cam2 2200685   450 </leader>
        <datafield ind1=" " ind2=" " tag="AVA">
          <subfield code="b">BIB_TEST_1</subfield>
          <subfield code="j">LOC_TEST_1</subfield>
        </datafield>
        <datafield ind1=" " ind2=" " tag="AVA">
          <subfield code="b">BIB_TEST_2</subfield>
          <subfield code="j">LOC_TEST_2</subfield>
        </datafield>
      </record>
    </bib>
  </bibs>`)

var fakeIln2rcr = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="iln2rcr">
  <query>
    <iln>ILN01</iln>
    <result>
      <library>
        <rcr>rcr000001</rcr>
        <shortname>BIB_TEST_1</shortname>
      </library>
      <library>
        <rcr>rcr000002</rcr>
        <shortname>BIB_TEST_2</shortname>
      </library>
      <library>
        <rcr>rcr000003</rcr>
        <shortname>BIB_TEST_3</shortname>
      </library>
    </result>
  </query>
  <query>
    <iln>ILN02</iln>
    <result>
      <library>
        <rcr>rcr000004</rcr>
        <shortname>BIB_TEST_4</shortname>
      </library>
    </result>
  </query>
</sudoc>`)

var fakeIln2rcr_2 = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="iln2rcr">
  <query>
    <iln>ILN02</iln>
    <result>
      <library>
        <rcr>rcr000004</rcr>
        <shortname>BIB_TEST_4</shortname>
      </library>
    </result>
  </query>
</sudoc>`)

var fakeIln2rcrError = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<sudoc service="iln2rcr">
<error>Nil</error>
</sudoc>`)

var marcAax = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<record>
  <leader>     cam0 22        450 </leader>
  <controlfield tag="008">Aax3</controlfield>
  <datafield tag="010" ind1=" " ind2=" ">
    <subfield code="a">978-2-253-02983-0</subfield>
  </datafield>
</record>`)

var marcOax = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<record>
  <leader>     cam0 22        450 </leader>
  <controlfield tag="008">Oax3</controlfield>
  <datafield tag="010" ind1=" " ind2=" ">
    <subfield code="a">978-2-253-02983-0</subfield>
  </datafield>
</record>`)

var marcError = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<error>Les données bibliographiques sont indéfinies
  <ppn>155075380</ppn>
</error>`)
